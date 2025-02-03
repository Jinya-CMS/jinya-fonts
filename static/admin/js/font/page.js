import { Alpine } from '../../../lib/alpine.js';
import { localize } from '../../../lib/jinya-alpine-tools.js';
import { get, httpDelete } from '../../../lib/jinya-http.js';
import confirm from '../../lib/ui/confirm.js';

import '../../lib/ui/toolbar-editor.js';

Alpine.data('fontsData', () => ({
  loading: true,
  createFontOpen: false,
  activeFilter: {
    name: '',
    category: 'All',
  },
  newFont: {
    name: '',
    category: 'Sans Serif',
    license: '',
    description: '',
  },
  get activeSideItem() {
    return this.$router.params.page;
  },
  fonts: [],
  get noFontsMessage() {
    switch (this.activeSideItem) {
      case 'google':
        return localize({ key: 'font-list-no-google' });
      case 'custom':
        return localize({ key: 'font-list-no-custom' });
      default:
        return localize({ key: 'font-list-no-all' });
    }
  },
  get filteredFonts() {
    let filteredFonts = this.fonts.filter((f) => f.name.toLowerCase().includes(this.activeFilter.name.toLowerCase()));
    if (this.activeFilter.category !== 'All') {
      filteredFonts = filteredFonts.filter((f) => f.category === this.activeFilter.category);
    }
    return filteredFonts;
  },
  get title() {
    switch (this.activeSideItem) {
      case 'google':
        return localize({ key: 'font-list-google-fonts' });
      case 'custom':
        return localize({ key: 'font-list-custom-fonts' });
      default:
        return localize({ key: 'font-list-all-fonts' });
    }
  },
  async init() {
    this.fonts = await get('/api/admin/font/' + this.activeSideItem);
    this.loading = false;
  },
  designersToString(designers) {
    return designers?.map((d) => d.name).join(', ') ?? '';
  },
  openCreateFont() {
    this.newFont.name = '';
    this.newFont.category = 'Sans Serif';
    this.newFont.license = '';
    this.newFont.description = '';
    this.createFontOpen = true;
  },
  syncFonts() {
    Alpine.store('syncProgress').start();
  },
  async deleteFont(font) {
    if (
      await confirm({
        title: localize({ key: 'delete-font-title' }),
        approveLabel: localize({ key: 'delete-font-confirm-label' }),
        declineLabel: localize({ key: 'delete-font-decline-label' }),
        message: localize({ key: 'delete-font-message', values: font }),
        negative: true,
      })
    ) {
      await httpDelete(`/api/admin/font/${font.name}`);
      this.fonts = this.fonts.filter((f) => f.name !== font.name);
    }
  },
  createFont() {},
}));
