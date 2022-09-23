# Release button

Make ArgoCD releases with physical button hooked up to Raspberry Pi.

## TODO
- [X] Argo API
- [X] Wait for button clicks
- [X] Make release on click
- [X] Fetch statuses of apps (which needs syncing)
- [X] Ignore some projects
- [X] Show status with LED
---
Future:
- [ ] Safeguards to prevent accidental releases
- [ ] Hoop up a screen



mockgen -source=internal/argoApi/models.go -destination=internal/mocks/argoargoApi.go -package=mocks IArgoargoApi
