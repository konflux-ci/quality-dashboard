### In progress

Dex-config example

```yaml
# The base path of dex and the external name of the OpenID Connect service.
# This is the canonical URL that all clients MUST use to refer to dex. If a
# path is provided, dex’s HTTP service will listen at a non-root URL.
issuer: http://127.0.0.1:5556/dex # remember to change this to RHTAP url when deploy
# The storage configuration determines where dex stores its state. Supported
# options include SQL flavors and Kubernetes third party resources.
#
# See the documentation (https://dexidp.io/docs/storage/) for further information.
storage:
  type: sqlite3
web:
  http: 0.0.0.0:5556
  allowedOrigins: ["*"]
oauth2:
  skipApprovalScreen: true
telemetry:
  http: 0.0.0.0:5558
  enableProfiling: true
expiry:
  deviceRequests: "5m"
  signingKeys: "6h"
  idTokens: "1h"
  refreshTokens:
    reuseInterval: "3s"
    validIfNotUsedFor: "2160h"
    absoluteLifetime: "3960h"
staticClients:
- id: redhat-quality-studio-app
  redirectURIs:
  - "http://localhost:9000/home/overview"
  - "http://localhost:9000/login"
  name: "Red Hat Quality Studio"
  public: true
connectors:
- type: github
  id: github
  name: GitHub
  config:
    orgs:
    - name: redhat-appstudio-qe
      teams:
      - admins
    clientID: <your github client id>
    clientSecret: <your github client secret
    redirectURI: http://127.0.0.1:5556/dex/callback # remember to change this to RHTAP url when deploy
```