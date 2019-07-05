package plugin

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// The plugin type downloaded from a Git repository.
type Git struct {
	// Remote URL
	URL string
	// The required revision. This can be a branch, a tag or a commit hash.
	requiredRevision string
	// We reuse the `Dir` plugin type to load the plugin into zsh.
	Dir        Dir
	repository *git.Repository
	// The target commit hash set by `CheckUpdate`.
	update *plumbing.Hash
}

// NewGit creates a new Git plugin. URL must come without a scheme (like
// github.com/robbyrussel/oh-my-zsh).
func NewGit(URL string, requiredRevision string, root string) Git {
	return Git{
		URL:              URL,
		requiredRevision: requiredRevision,
		Dir:              Dir{Path: filepath.Join(root, URL)},
	}
}

func (p *Git) Load() ([]string, []string, error) {
	return p.Dir.Load()
}

func (p *Git) CheckUpdate() (message *string, err error) {
	p.repository, err = git.PlainOpen(p.Dir.Path)
	if err == git.ErrRepositoryNotExists {
		return nil, NotInstalled
	} else if err != nil {
		return nil, errors.Wrap(err, "while opening the repository")
	}

	currentHead, err := p.repository.Head()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read repository HEAD")
	}
	if currentHead == nil {
		return nil, errors.New("cannot read repository HEAD")
	}

	currentVersion := currentHead.Hash()

	fetchOptions := git.FetchOptions{}
	if err := fetchOptions.Validate(); err != nil {
		return nil, errors.Wrap(err, "while fetching the repository")
	}
	if err := p.repository.Fetch(&fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, errors.Wrap(err, "while fetching the repository")
	}

	// because we fetch, not pull, we need to check the remote branches
	newVersionRemote := plumbing.NewRemoteReferenceName("origin", p.requiredRevision)
	newVersion, err := p.repository.ResolveRevision(plumbing.Revision(newVersionRemote))
	if err != nil {
		newVersion, err = p.repository.ResolveRevision(plumbing.Revision(p.requiredRevision))
		if err != nil {
			newVersionLocal := plumbing.NewHash(p.requiredRevision)
			newVersion = &newVersionLocal
			if o, _ := p.repository.CommitObject(newVersionLocal); o == nil {
				return nil, errors.New("failed to get the revision")
			}
			return nil, errors.New("failed to get the revision")
		}
	}

	if *newVersion == currentVersion {
		return nil, UpToDate
	}

	updateString := fmt.Sprintf(
		"%s: update from %s to %s",
		p.requiredRevision,
		currentVersion.String()[:7],
		newVersion.String()[:7],
	)
	p.update = newVersion

	return &updateString, nil
}

func (p *Git) InstallUpdate() error {
	// install if an existing installation not found
	if p.repository == nil {
		parentPath := filepath.Dir(p.Dir.Path)
		if err := os.MkdirAll(parentPath, os.ModePerm); err != nil && !os.IsExist(err) {
			return errors.Wrap(err, "while creating github plugin object")
		}

		repositoryURL := fmt.Sprintf("https://%s.git", p.URL)
		cloneOptions := git.CloneOptions{URL: repositoryURL}
		if _, err := git.PlainClone(p.Dir.Path, false, &cloneOptions); err != nil {
			return errors.Wrap(err, "while cloning the repository")
		}
		return nil
	}

	if p.update == nil {
		return errors.New("no update available")
	}

	worktree, err := p.repository.Worktree()
	if err != nil {
		return errors.Wrap(err, "checkout error")
	}

	return worktree.Checkout(&git.CheckoutOptions{Hash: *p.update})
}
