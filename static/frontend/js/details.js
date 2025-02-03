import Alpine from '/static/lib/alpine.js';

Alpine.data('detailData', () => ({
  previewType: 'lorem',
  customPreviewText: null,
  previewSize: 24,
  styles: [],
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
  get fontLink() {
    const styles = this.styles.sort((a, b) => a.localeCompare(b)).map((style) => style.split('-'));
    if (styles.length > 0) {
      const prefix = styles.filter(([, style]) => style !== 'normal').length > 0 ? 'wght,ital' : 'wght';
      let weightAndStyles = '';
      if (prefix === 'wght') {
        weightAndStyles = styles.map(([weight]) => `${weight}`).join(';');
      } else {
        weightAndStyles = styles.map(([weight, style]) => `${weight},${style}`).join(';');
      }

      return `${location.origin}/css2?family=${window.fontName}:${prefix}@${weightAndStyles}`;
    }

    return `${location.origin}/css2?family=${window.fontName}`;
  },
  get htmlCode() {
    return `<link rel="stylesheet" href="${this.fontLink}">`;
  },
  get fontCss() {
    return `@import url('${this.fontLink}');`;
  },
}));

Alpine.start();
