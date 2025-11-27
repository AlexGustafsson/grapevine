Grapevine dead simple app-less notification system for iOS using declerative
web push.

Use cases: Grafana (web push). Support whatever Gotify supports.

Could use E2EE by including public key when sending web push subscription to
server, then include encrypted message in data. Would require service worker in
the app.

Somehow allow selecting favicon / hosting multiple ingresses. That way, when
adding the app to the home screen, the name and logo can be chosen, which is
then shown as "getting a notification from Grafana". Could be done simply by
paths (if webserver sees /grafana, server grafana icon as favicon), or by
browser fingerprinting (select grafana in Safari, then when adding to screen
the webserver serves the icon from manifest).

iOS (26.1) requests the image when adding opening to add it to the device:

- `/assets/favicon-DRNirk6r.png` (from `index.html`)
- `/assets/icon.png` (from `manifest.json`)
- `/apple-touch-icon-120x120-precomposed.png`
- `/apple-touch-icon-120x120.png`
- `/apple-touch-icon.png`
- `/apple-touch-icon.png`

Once added, it's unclear when it's updated again.

User-Agent for the `/assets/favicon` one is
`NetworkingExtension/8622.2.11.10.8 Network/5569.42.2 iOS/26.1`, for the others
it's `SafariViewService/8622.2.11.10.8 CFNetwork/3860.200.71 Darwin/25.1.0`,
so the user agent is probably not there? No referer, so we can't use that.

The `manifest.json` file is fetched with a referer, could maybe be used - but
wouldn't be persistent. Fetched with Safari's user agent string. The manifest
file is fetched before opening to add the web app, so likely cannot be used.

The `manifest.json` is fetched immediately after `index.html`.

The following should work:

User goes to home page. Selects "Grafana" (pre-configured I guess). Taken to a
new site - with hard reload! The manifest.json link is updated with one
pre-configured / templated for Grafana instead, with unique URLs and start url?

- Alertmanager: <https://github.com/prometheus/alertmanager/blob/main/api/v2/openapi.yaml>
- Grafana: <https://grafana.com/docs/grafana/latest/alerting/configure-notifications/manage-contact-points/integrations/webhook-notifier/>
