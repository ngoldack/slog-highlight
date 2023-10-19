# ci

This package contains the dagger code for the CI/CD pipeline.

## Why separate package?

The CI/CD pipeline is a separate package because it is not required for the library to function. It is only required for the development of the library. To prevent the CI/CD pipeline from being included in the library, it is placed in a separate package.
