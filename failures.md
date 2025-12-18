# Automation API related failures

## E2E tests:
Small fix to do (YAML export => flow execution case):
- -- FAIL: chainsaw/v4-export-import (14.07s)
Skipped tests as need Automation API support for notifications:
- -- FAIL: chainsaw/remove-notification (69.51s)
- -- FAIL: chainsaw/update-notification-events (38.17s)
- -- FAIL: chainsaw/update-notification-grouprefs (50.52s)

## Integration tests:
Skipped tests as need Automation API support for notifications:
* [FAIL] Create [It] should add notification to API after creation [integration, withContext]
      `/Users/benoit.bordigoni/src/gravitee-kubernetes-operator/test/integration/apidefinition/v4/create_withContext_andNotificationWithGroupRef_test.go:99`
