import { Alpine } from '../../../lib/alpine.js';
import { get, httpDelete } from '../../../lib/jinya-http.js';

import '../../lib/ui/toolbar-editor.js';
import confirm from '../../lib/ui/confirm.js';
import { localize } from '../../../lib/jinya-alpine-tools.js';

Alpine.data('fontDetailsData', () => ({
  loading: true,
  activeSideItem: 'details',
  font: null,
  activeDesigner: null,
  editDesignerOn: false,
  get fontName() {
    return this.$router.params.name;
  },
  async init() {
    this.font = await get(`/api/admin/font/${this.fontName}`);
    this.loading = false;
    if (this.font.designers.length > 0) {
      this.activeDesigner = this.font.designers[0];
    }
  },
  changeSideItem(item) {
    this.activeSideItem = item;
  },
  openEditFontDialog() {},
  async deleteFont() {
    if (
      await confirm({
        title: localize({ key: 'delete-font-title' }),
        approveLabel: localize({ key: 'delete-font-confirm-label' }),
        declineLabel: localize({ key: 'delete-font-decline-label' }),
        message: localize({ key: 'delete-font-message', values: this.font }),
        negative: true,
      })
    ) {
      await httpDelete(`/api/admin/font/${this.font.name}`);
      this.$router.navigate('/font/all');
    }
  },
}));
