import Alpine from '/static/lib/alpine.js';
import { get } from '../../lib/jinya-http.js';

Alpine.data('indexData', () => ({
  loading: true,
  filterText: '',
  previewType: 'lorem',
  customPreviewText: null,
  previewSize: 24,
  sansEnabled: true,
  serifEnabled: true,
  handwrittenEnabled: true,
  displayEnabled: true,
  monospaceEnabled: true,
  fonts: [],
  get filteredFonts() {
    return this.fonts.filter((font) => {
      let nameMatch = font.name.toLowerCase().includes(this.filterText.toLowerCase());
      let categoryMatch =
        (font.category === 'Sans Serif' && this.sansEnabled) ||
        (font.category === 'Serif' && this.serifEnabled) ||
        (font.category === 'Handwritten' && this.handwrittenEnabled) ||
        (font.category === 'Display' && this.displayEnabled) ||
        (font.category === 'Monospace' && this.monospaceEnabled);

      return nameMatch && categoryMatch;
    });
  },
  get previewText() {
    switch (this.previewType) {
      case 'alphabet':
        return 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
      case 'lorem':
        return 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.';
      case 'numbers':
        return '12345678890';
    }

    return this.customPreviewText;
  },
  async init() {
    this.fonts = await get('/api/font');
    this.loading = false;
    this.$watch('previewType', (previewType) => {
      if (previewType !== 'custom') {
        this.customPreviewText = '';
      }
    });
  },
  styleUrl(font) {
    return `/css2?family=${font.name}`;
  },
  designer(font) {
    return font.designers.map((d) => d.name).join(', ');
  },
}));

Alpine.start();
