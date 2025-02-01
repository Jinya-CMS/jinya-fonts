import { Alpine } from '../../../lib/alpine.js';
import { get, httpDelete, post, put } from '../../../lib/jinya-http.js';

import '../../lib/ui/toolbar-editor.js';
import confirm from '../../lib/ui/confirm.js';
import { localize } from '../../../lib/jinya-alpine-tools.js';

Alpine.data('fontDetailsData', () => ({
  loading: true,
  activeSideItem: 'details',
  font: null,
  activeDesigner: null,
  editDesignerOn: false,
  editFontOpen: false,
  newDesigner: {
    name: '',
    bio: '',
  },
  newFileOpen: false,
  newFile: {
    file: undefined,
    weight: '400',
    style: 'normal',
  },
  editFileOpen: false,
  editFile: {
    file: undefined,
    weight: '400',
    style: 'normal',
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
    this.font.editDescription = this.font.description;
    this.font.editLicense = this.font.license;
    this.font.editCategory = this.font.category;
    this.editFontOpen = true;
  },
  openNewFileDialog() {
    this.newFile = {
      file: undefined,
      weight: '400',
      style: 'normal',
    };
    this.newFileOpen = true;
  },
  openEditFileDialog(file) {
    this.editFile = {
      file: undefined,
      weight: file.weight,
      style: file.style,
    };
    this.editFileOpen = true;
  },
  editDesigner() {
    this.activeDesigner.editBio = this.activeDesigner.bio;
    this.editDesignerOn = true;
  },
  async updateFont() {
    this.font.description = this.font.editDescription;
    this.font.license = this.font.editLicense;
    this.font.category = this.font.editCategory;
    await put(`/api/admin/font/${this.fontName}`, {
      description: this.font.editDescription,
      license: this.font.editLicense,
      category: this.font.editCategory,
    });
    this.editFontOpen = false;
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
    this.font.designers = [...this.font.designers.filter((d) => d.name !== this.activeDesigner.name), updatedDesigner];
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
      this.font.designers = this.font.designers.filter((d) => d.name !== this.activeDesigner.name);
      this.activeDesigner = this.font.designers.at(0) ?? null;
    }
  },
  async deleteFile(file) {
    if (
      await confirm({
        title: localize({ key: 'delete-file-title' }),
        approveLabel: localize({ key: 'delete-file-confirm-label' }),
        declineLabel: localize({ key: 'delete-file-decline-label' }),
        message: localize({
          key: 'delete-file-message',
          values: {
            type: file.type,
            weight: localize({ key: `font-weight-${file.weight}` }),
            style: localize({ key: `font-style-${file.style}` }),
          },
        }),
        negative: true,
      })
    ) {
      await httpDelete(`/api/admin/font/${this.fontName}/file/${file.weight}.${file.style}.${file.type}`);
      this.font.fonts = this.font.fonts.filter((f) => f.path !== file.path);
    }
  },
  async createNewFile() {
    const type = this.newFile.file.name.split('.').reverse()[0];
    await post(`/api/admin/font/${this.fontName}/file/${this.newFile.weight}.${this.newFile.style}.${type}`, this.newFile.file);
    this.font.fonts = await get(`/api/admin/font/${this.fontName}/file`);
    this.newFileOpen = false;
  },
  async updateFile() {
    const type = this.editFile.file.name.split('.').reverse()[0];
    await post(`/api/admin/font/${this.fontName}/file/${this.editFile.weight}.${this.editFile.style}.${type}`, this.editFile.file);
    this.font.fonts = await get(`/api/admin/font/${this.fontName}/file`);
    this.editFileOpen = false;
  },
}));
