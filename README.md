# Release button

Make ArgoCD releases with physical button hooked up to Raspberry Pi.

## TODO
- [X] Argo API
- [X] Wait for button clicks
- [X] Make release on click
- [X] Fetch statuses of apps (which needs syncing)
- [X] Ignore some projects
- [X] Show status with LED
- [X] Safeguards to prevent accidental releases
---
Future:
- [ ] Update refreshInterval wihout restart
- [X] Localhost in web ui is not correct 


## Notes

To generate test mock for an interface:
`mockgen -source=internal/argoApi/models.go -destination=internal/mocks/argoApi.go -package=mocks IArgoApi`
