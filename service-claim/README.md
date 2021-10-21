# Service Claim

## Development

This project was scaffolded with:

``` bash
$ brew install operator-sdk
$ operator-sdk init --domain vmware.com --owner "Guido Westenberg, VMware Inc" --project-name service-claim --repo github.com/cf-platform-eng/service-claim
$ operator-sdk create api --group tsmgr --kind ServiceClaim --namespaced false --version v1 --resource --controller
```
