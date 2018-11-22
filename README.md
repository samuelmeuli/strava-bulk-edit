# Strava Bulk Edit

Unfortunately, Strava doesn't allow users to update information about more than one activity at once. This small command-line utility saves you the tedious work of updating every single activity by hand.


## Usage

Simply download the binary from the [releases page](https://github.com/samuelmeuli/strava-bulk-edit/releases/latest) and execute it using the following command:

```sh
./strava-bulk-edit --help
```


## Example

For example, you could make all of your Strava activities before 2018 private using the following command:

```sh
./strava-bulk-edit visibility only_me --to 2018-01-01
```


## Options

The following activity attributes can be updated using this command-line script:

* title
* description
* commute
* type
* visibility

You can update all of your activities using the `--all` flag or activities in a date range using the `--from` and/or `--to` flags.

Run `./strava-bulk-edit --help` for detailed descriptions of the commands and information about allowed values.
