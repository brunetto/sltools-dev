This is the dev version of [`sltools`](https://github.com/brunetto/sltools). Please refer to that.    

### Management

To pull changes from `sltools` stable:

```bash
git pull upstream
```

To pull or push changes from/to `sltools-dev`

```bash
git pull/push
```
or

```bash
git pull/push origin master
```

### History of what I've done

```bash
git clone git@github.com:brunetto/sltools.git					# clone the master branch of the stable repo here
git remote rm origin											# remove old (stable) origin
git remote add origin git@github.com:brunetto/sltools-dev.git	# add new origin
git remote add upstream git@github.com:brunetto/sltools.git		# add stable repo as `upstream`
git pull --force upstream dev									# download the stable repo dev branch and merge with the local master one
git pull upstream master										# try to update with new changes in the stable/master
git push --set-upstream origin master							# push to the dev repo master branch

```
