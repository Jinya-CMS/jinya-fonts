import { setup } from '/static/lib/jinya-alpine-tools.js';

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
    languages: { de: MessagesDe, en: MessagesEn }
  });
});