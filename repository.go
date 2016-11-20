package subversion

import (
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
)

// IsRepository checks if a file represents a
// subversion repository.
func IsRepository(path string) bool {
	cmd := exec.Command("svnlook", "uuid", path)
	if err := cmd.Start(); err != nil {
		return false
	}

	err := cmd.Wait()
	if err != nil {
		return false
	}

	return true
}

// Backup creates or updates a backup for a repository.
func Backup(repo, dest string) error {
	if IsRepository(dest) {
		return HotCopy(repo, dest, true)
	}
	return HotCopy(repo, dest, false)
}

// HotCopy creates a copy of the repository to the
// new destination.
func HotCopy(repo, dest string, incremental bool) error {
	var cmd *exec.Cmd
	if incremental {
		cmd = exec.Command("svnadmin", "hotcopy", "--incremental", repo, dest)
	} else {
		cmd = exec.Command("svnadmin", "hotcopy", repo, dest)
	}

	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "Error running svnadmin hotcopy")
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "svnadmin hotcopy failed")
	}

	return nil
}

// LatestRevision returns the most recent (youngest)
// revision number for the given repostory
func LatestRevision(repo string) (int, error) {
	cmd := exec.Command("svnlook", "youngest", repo)
	out, err := cleanCmdOutput(cmd)
	if err != nil {
		return 0, errors.Wrap(err, "svnlook youngest failed")
	}

	rev, err := strconv.Atoi(string(out))
	if err != nil {
		return 0, errors.Wrap(err, "Unable to convert svnlook youngest output to a revision number")
	}

	return rev, nil
}

// UUID returns the repository's universal unique identifier.
func UUID(repo string) (string, error) {
	cmd := exec.Command("svnlook", "uuid", repo)
	out, err := cleanCmdOutput(cmd)
	if err != nil {
		return "", errors.Wrap(err, "svnlook uuid failed")
	}

	return string(out), err
}
