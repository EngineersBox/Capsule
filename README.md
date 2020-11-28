# Capsule

Ship shipping shipping ship shipping shipping ships with containers. Aka Docker clone.

## GoLang Environment

In order to have correct building, Capsule will need to be build on a Linux based system as it uses native Linux kernel calls.

If you are however, developing on a non-Linux OS and want to have correct linting and build/run response then you will need to configure the `GOOS` environment variable.

To do this, simply run the following command:

```bash
go env -w GOOS=linux
```
