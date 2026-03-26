## 1. Rename Command to `install`

- [x] 1.1 Change `Use` field from `"init"` to `"install [--path <dir>]"`
- [x] 1.2 Update the variable name `initCmd` to `installCmd` throughout `cmd/init.go`
- [x] 1.3 Update `rootCmd.AddCommand` call to register `installCmd`

## 2. Update Help Text

- [x] 2.1 Replace `Short` description with a one-line summary (e.g., `Install git-syn into a repository`)
- [x] 2.2 Replace `Long` description explaining that the command installs a pre-push hook and registers remotes for synchronization

## 3. Add --path Flag

- [x] 3.1 Declare a `path` string variable in the `init()` function
- [x] 3.2 Register `--path` as a local string flag on `installCmd` with empty string default
- [x] 3.3 In the `Run` function, resolve the target path: use `--path` value if non-empty, otherwise fall back to `os.Getwd()`

## 4. Define Embedded Hook Template

- [x] 4.1 Declare a `hookContent` string constant containing the hook shell script with `{{Command}}` as the placeholder, following the Git LFS `hookBaseContent` pattern:
  - Check `command -v git-syn`; if absent print a warning to stderr and exit 2
  - If present, delegate to `git syn {{Command}} "$@"`

## 5. Implement `write_hook`

- [x] 5.1 Define `write_hook(path, content string) error`
- [x] 5.2 Call `os.WriteFile(path, []byte(content+"\n"), 0600)` and return any error

## 6. Implement `set_file_permissions`

- [x] 6.1 Define `set_file_permissions(path string, mode os.FileMode) error`
- [x] 6.2 Call `os.Chmod(path, mode)` and return any error

## 7. Implement `install_hook`

- [x] 7.1 Define `install_hook(hookType, hooksDir string) error`
- [x] 7.2 Replace `{{Command}}` in `hookContent` with `hookType`
- [x] 7.3 Call `write_hook(filepath.Join(hooksDir, hookType), content)`; return error on failure
- [x] 7.4 Call `set_file_permissions(filepath.Join(hooksDir, hookType), 0700)`; return error on failure

## 8. Implement `init_repo`

- [x] 8.1 Define `init_repo(path string)`
- [x] 8.2 Call `git.PlainOpen(path)`; log fatal with a descriptive message if it is not a git repository
- [x] 8.3 Iterate remotes; accept `https://` and `git@` URLs; warn and skip others
- [x] 8.4 Write accepted remotes to `<path>/.gitremotes` in gitconfig format
- [x] 8.5 Read `.gitremotes` back and parse each `[remote "name"]` / `url` entry
- [x] 8.6 Load repo config via `repo.Config()`; add each parsed remote if not already present; save via `repo.SetConfig()`
- [x] 8.7 Handle all errors with descriptive `log.Fatal` messages

## 9. Wire Up Run Function

- [x] 9.1 Resolve target path (task 3.3)
- [x] 9.2 Call `install_hook("pre-push", filepath.Join(targetPath, ".git", "hooks"))`
- [x] 9.3 Call `init_repo(targetPath)`
- [x] 9.4 Print `Updated git hooks. Git SYN initialized.`
