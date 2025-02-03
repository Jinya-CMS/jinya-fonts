import { localize, setup } from '/static/lib/jinya-alpine-tools.js';
import { Alpine } from '../lib/alpine.js';
import { post } from '../lib/jinya-http.js';

document.addEventListener('DOMContentLoaded', async () => {
  const MessagesDe = await (await fetch('/static/admin/langs/messages.de.json')).json();
  const MessagesEn = await (await fetch('/static/admin/langs/messages.en.json')).json();

  await setup({
    defaultArea: 'font',
    defaultPage: 'page',
    baseScriptPath: '/static/admin/js/',
    routerBasePath: '/admin',
    openIdClientId: window.jinyaConfig.openIdClientId,
    openIdUrl: window.jinyaConfig.openIdUrl,
    openIdCallbackUrl: window.jinyaConfig.openIdCallbackUrl,
    languages: { de: MessagesDe, en: MessagesEn },
    afterSetup() {
      Alpine.store('syncProgress', {
        fontsSyncSuccess: false,
        fontsSyncFailure: false,
        fontsSyncing: false,
        get syncStatusText() {
          if (this.fontsSyncing) {
            return localize({ key: 'synchronizing-google-fonts' });
          } else if (this.fontsSyncSuccess) {
            return localize({ key: 'synchronization-success' });
          } else if (this.fontsSyncFailure) {
            return localize({ key: 'synchronization-error' });
          }

          return '';
        },
        start() {
          this.fontsSyncing = true;
          post('/api/admin/font/google')
            .then((result) => {
              this.fontsSyncSuccess = true;
            })
            .catch((err) => {
              this.fontsSyncFailure = true;
            })
            .finally(() => {
              this.fontsSyncing = false;
            });
        },
      });
      Alpine.directive('boolean-display', (el, { expression, modifiers }, { Alpine, effect }) => {
        effect(() => {
          if (Alpine.evaluate(el, expression)) {
            el.innerHTML = `<svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M20 6 9 17l-5-5" />
                </svg>`;
          } else {
            el.innerHTML = `<svg
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M18 6 6 18" />
                  <path d="m6 6 12 12" />
                </svg>`;
          }
        });
      });
      Alpine.directive('font-license', (el, { expression, modifiers }, { Alpine, effect }) => {
        effect(() => {
          const license = Alpine.evaluate(el, expression);
          switch (license) {
            case 'apache2':
              el.innerHTML =
                '<a target="_blank" href="https://www.apache.org/licenses/LICENSE-2.0">Apache License, Version 2.0</a>';
              break;
            case 'ofl':
              el.innerHTML = '<a target="_blank" href="https://openfontlicense.org/">Open Font License</a>';
              break;
            case 'ufl':
              el.innerHTML = '<a target="_blank" href="https://font.ubuntu.com/ufl/">Ubuntu Font License</a>';
              break;
            default:
              el.innerHTML = license;
              break;
          }
        });
      });
    },
  });
});
