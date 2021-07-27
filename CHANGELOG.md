## 0.4.0 (July 27, 2021)

IMPROVEMENTS:

- Update to the latest version of the Terraform Plugin SDK. This should improve compatibility with Terraform v1.0 ([#16](https://github.com/liatrio/terraform-provider-harbor/pull/16))

BUG FIXES:

- Fix an issue with deleting the `harbor_project` resource when chartmuseum is disabled ([#16](https://github.com/liatrio/terraform-provider-harbor/pull/16))

## 0.3.3 (March 3, 2021)

BUG FIXES:

- Fixes errors with `harbor_project` resources that were either imported or generated with an earlier provider version, and had never set the `auto_scan` attribute ([#14](https://github.com/liatrio/terraform-provider-harbor/pull/14))

## 0.3.2 (February 25, 2021)

FEATURES:

- Adds support for the `harbor_label` resource ([#13](https://github.com/liatrio/terraform-provider-harbor/pull/13))

## 0.3.1 (February 19, 2021)

FEATURES:

- Adds support for the `harbor_webhook` resource ([#12](https://github.com/liatrio/terraform-provider-harbor/pull/12))

IMPROVEMENTS:

- Adds the `auto_scan` attribute to the `harbor_project` resource ([#12](https://github.com/liatrio/terraform-provider-harbor/pull/12))

## 0.3.0 (February 16, 2021)

BREAKING CHANGES:

- Updates client library for use with Harbor 2.0, support for Harbor 1.0 is dropped ([#11](https://github.com/liatrio/terraform-provider-harbor/pull/11))

IMPROVEMENTS:

- Adds `expires_at` attribute to `harbor_robot_account` resource ([#11](https://github.com/liatrio/terraform-provider-harbor/pull/11))

## 0.2.0-pre (August 26, 2020)

IMPROVEMENTS:

- Updates CI and documentation for releasing to Terraform Registry ([#8](https://github.com/liatrio/terraform-provider-harbor/pull/8))

BUG FIXES:

- Resources now correctly recover from 404 errors from Harbor ([#7](https://github.com/liatrio/terraform-provider-harbor/pull/7))

## 0.1.0 (April 2, 2020)

BREAKING CHANGES:

- The `public` attribute in the `harbor_project` resource was changed from string to bool
- The `project_id` computed attribute in the `harbor_project` resource was removed, to be replaced by the `id` attribute.
- The `robot_account_access` attribute in the `harbor_robot_account` was renamed to `access`

IMPROVEMENTS:

- Adds validation for several attributes which previously would've failed with API errors

BUG FIXES:

- Deleting a `harbor_project` now works correctly when project contains respositories or helm charts ([#6](https://github.com/liatrio/terraform-provider-harbor/pull/6))

## 0.0.4 (March 27, 2020)

BUG FIXES:

- Syntax fixes for automated release

## 0.0.3 (March 27, 2020)

BUG FIXES:

- Authentication fixes for automated release

## 0.0.2 (March 27, 2020)

Initial Release
