import { Alpine } from '../../../lib/alpine.js';
import { get, put } from '../../../lib/jinya-http.js';

Alpine.data('configData', () => ({
  loading: false,
  enableSync: false,
  syncInterval: '0 0 * * *',
  onlyFetchFonts: [],
  addFilterName: '',
  removeFilteredFont(font) {
    this.onlyFetchFonts = this.onlyFetchFonts.filter((f) => f !== font);
  },
  addFilter() {
    this.onlyFetchFonts = [...this.onlyFetchFonts, this.addFilterName];
    this.addFilterName = '';
  },
  async loadData() {
    this.loading = true;
    const settings = await get('/api/admin/settings');
    this.onlyFetchFonts = settings.filterByName;
    this.syncInterval = settings.syncInterval;
    this.enableSync = settings.syncEnabled;
    this.loading = false;
  },
  async init() {
    await this.loadData();
  },
  async saveSettings() {
    await put(`/api/admin/settings`, {
      filterByName: this.onlyFetchFonts,
      syncEnabled: this.enableSync,
      syncInterval: this.syncInterval,
    });
  },
}));
