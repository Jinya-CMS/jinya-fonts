import { Alpine } from '../../../lib/alpine.js';
import { get, httpDelete, post } from '../../../lib/jinya-http.js';

import '../../lib/ui/toolbar-editor.js';
import confirm from '../../lib/ui/confirm.js';
import { localize } from '../../../lib/jinya-alpine-tools.js';

Alpine.data('fontDetailsData', () => ({
  loading: true,
  activeSideItem: 'details',
  font: null,
  activeDesigner: null,
  editDesignerOn: false,
  newDesigner: {
    name: '',
    bio: '',
  },
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
  openEditFontDialog() {
  },
  editDesigner() {
    this.activeDesigner.editBio = this.activeDesigner.bio;
    this.editDesignerOn = true;
  },
  async createNewDesigner() {
    const newDesigner = await post(`/api/admin/font/${this.fontName}/designer`, this.newDesigner);
    this.font.designers = [...this.font.designers, newDesigner];
    this.newDesigner.name = '';
    this.newDesigner.bio = '';
  },
  async updateDesigner() {
    await httpDelete(`/api/admin/font/${this.fontName}/designer/${this.activeDesigner.name}`);
    const updatedDesigner = await post(`/api/admin/font/${this.fontName}/designer`, {
      bio: this.activeDesigner.editBio,
      name: this.activeDesigner.name,
    });
    this.font.designers = [...this.font.designers.filter(d => d.name !== this.activeDesigner.name), updatedDesigner];
    this.activeDesigner = updatedDesigner;
    this.activeDesigner.editBio = updatedDesigner.bio;
    this.editDesignerOn = false;
  },
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
  async removeDesigner() {
    if (
      await confirm({
        title: localize({ key: 'delete-designer-title' }),
        approveLabel: localize({ key: 'delete-designer-confirm-label' }),
        declineLabel: localize({ key: 'delete-designer-decline-label' }),
        message: localize({ key: 'delete-designer-message', values: this.activeDesigner }),
        negative: true,
      })
    ) {
      await httpDelete(`/api/admin/font/${this.fontName}/designer/${this.activeDesigner.name}`);
      this.font.designers = this.font.designers.filter(d => d.name !== this.activeDesigner.name);
      this.activeDesigner = this.font.designers.at(0) ?? null;
    }
  },
}));
