package plugin

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
)

type GitHub struct {
	username        string
	repositoryName  string
	requiredVersion string
	root            string
	Dir             *Dir
	repository      *git.Repository
	update          *plumbing.Hash
}

func MakeGitHub(root string, params []string) (*Plugin, error) {
	if len(params) != 4 {
		return nil, errors.New("invalid number of parameters")
	}

	requiredRevision := "master"
	if params[3] != "" {
		requiredRevision = params[3]
	}

	github, err := NewGitHub(params[0], params[1], requiredRevision, root)
	plugin := Plugin(github)

	return &plugin, err
}

func NewGitHub(
	username string,
	repository string,
	requiredVersion string,
	root string,
) (*GitHub, error) {
	var dir *Dir

	path := filepath.Join(root, "Plugins", "github.com", username, repository)
	stat, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "while creating github plugin object")
	} else if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			return nil, errors.New("the provided path is not a directory: " + path)
		}
		dir = &Dir{Path: path}
	}

	ret := GitHub{
		username:        username,
		repositoryName:  repository,
		requiredVersion: requiredVersion,
		root:            root,
		Dir:             dir,
	}

	return &ret, nil
}

func (p *GitHub) Load() ([]string, []string, error) {
	if p.Dir == nil {
		return nil, nil, errors.New("plugin is not present on the drive")
	}
	return p.Dir.Load()
}

func (p *GitHub) clone() error {
	parentPath := filepath.Join(p.root, "Plugins", "github.com", p.username)
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "while creating github plugin object")
	}

	path := filepath.Join(p.root, "Plugins", "github.com", p.username, p.repositoryName)

	repositoryURL := fmt.Sprintf("https://github.com/%s/%s.git", p.username, p.repositoryName)
	cloneOptions := git.CloneOptions{URL: repositoryURL}
	if _, err := git.PlainClone(path, false, &cloneOptions); err != nil {
		return errors.Wrap(err, "while cloning the repository")
	}

	p.Dir = &Dir{Path: path}

	return nil
}

func (p *GitHub) CheckUpdate() (*string, error) {
	if p.Dir == nil {
		return nil, NotInstalled
	}

	repo, err := git.PlainOpen(p.Dir.Path)
	if err != nil {
		return nil, err
	}

	currentHead, err := repo.Head()
	if err != nil {
		return nil, errors.Wrap(err, "cannot read chain head")
	}
	if currentHead == nil {
		return nil, errors.New("cannot read chain head")
	}

	currentVersion := currentHead.Hash()

	fetchOptions := git.FetchOptions{}
	if err := fetchOptions.Validate(); err != nil {
		return nil, errors.Wrap(err, "while fetching the repositoryName")
	}
	if err := repo.Fetch(&fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, errors.Wrap(err, "while fetching the repositoryName")
	}

	newVersionRemote := plumbing.NewRemoteReferenceName("origin", p.requiredVersion)
	newVersion, err := repo.ResolveRevision(plumbing.Revision(newVersionRemote))
	if err != nil {
		newVersion, err = repo.ResolveRevision(plumbing.Revision(p.requiredVersion))
		if err != nil {
			newVersionLocal := plumbing.NewHash(p.requiredVersion)
			newVersion = &newVersionLocal
			if o, _ := repo.CommitObject(newVersionLocal); o == nil {
				return nil, errors.New("failed to get the revision")
			}
			return nil, errors.New("failed to get the revision")
		}
	}

	if *newVersion == currentVersion {
		return nil, nil
	}

	updateString := fmt.Sprintf(
		"%s: update from %s to %s",
		p.requiredVersion,
		currentVersion,
		newVersion,
	)

	p.update = newVersion
	p.repository = repo

	return &updateString, nil
}

func (p *GitHub) InstallUpdate() error {
	if p.Dir == nil {
		return p.clone()
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

func (p GitHub) GetPath() *string {
	if p.Dir != nil {
		return &p.Dir.Path
	}
	return nil
}
