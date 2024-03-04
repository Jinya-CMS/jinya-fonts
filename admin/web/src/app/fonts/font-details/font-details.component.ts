import { Component, Input, OnInit } from '@angular/core';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { Check, X } from 'lucide-angular';
import { Designer } from '../../api/models/designer';
import { Router } from '@angular/router';
import { FormControl, FormGroup } from '@angular/forms';
import { Metadata } from '../../api/models/metadata';

enum ActiveSideListItem {
  Details,
  Designers,
  Files
}

@Component({
  selector: 'app-font-details',
  templateUrl: './font-details.component.html',
  styleUrl: './font-details.component.scss'
})
export class FontDetailsComponent implements OnInit {
  @Input() fontName!: string;
  @Input() section!: string;

  font!: Webfont;

  activeDesigner: Designer | null = null;
  activeFile: Metadata | null = null;

  loading = true;
  loaded = false;
  deleteFontOpen = false;
  editDesignerOn = false;
  editFileOpen = false;
  removeDesignerOpen = false;
  removeFileOpen = false;

  activeSideListItem = ActiveSideListItem.Details;

  addDesignerForm = new FormGroup({
    name: new FormControl('', { nonNullable: true }),
    bio: new FormControl('', { nonNullable: false })
  });
  editDesignerForm = new FormGroup({
    bio: new FormControl('', { nonNullable: false })
  });

  protected readonly ActiveSideListItem = ActiveSideListItem;
  protected readonly Check = Check;
  protected readonly X = X;

  constructor(
    private apiClient: ApiService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadData();
    switch (this.section.toLowerCase()) {
      case 'designer':
        this.activeSideListItem = ActiveSideListItem.Designers;
        break;
      case 'file':
        this.activeSideListItem = ActiveSideListItem.Files;
        break;
      default:
        this.activeSideListItem = ActiveSideListItem.Details;
        break;
    }
  }

  changeSideItem($event: MouseEvent, item: ActiveSideListItem) {
    $event.preventDefault();
    this.activeSideListItem = item;
  }

  selectDesigner($event: Event | null, designer: Designer | null) {
    $event?.preventDefault();
    this.activeDesigner = designer;
    this.editDesignerForm.reset({ bio: designer?.bio });
    this.editDesignerOn = false;
  }

  deleteFont() {
    this.apiClient.deleteFontByName({ fontName: this.fontName }).subscribe(() => {
      this.router.navigateByUrl('/font');
    });
  }

  loadData(activeDesigner: Designer | null = null) {
    if (!this.loaded) {
      this.loading = true;
    }

    this.apiClient.getFontByName({ fontName: this.fontName }).subscribe((value) => {
      this.font = value;
      this.loading = false;
      if (activeDesigner) {
        this.selectDesigner(null, activeDesigner);
      } else {
        this.selectDesigner(null, this.font.designers?.at(0) ?? null);
      }
      this.addDesignerForm.reset();
      this.loaded = true;
      this.editFileOpen = false;
      this.removeDesignerOpen = false;
      this.removeFileOpen = false;
      this.activeFile = null;
    });
  }

  createNewDesigner() {
    if (this.addDesignerForm.invalid) {
      return;
    }

    this.apiClient
      .addFontDesigner({
        fontName: this.fontName,
        body: { name: this.addDesignerForm.get('name')?.value ?? '', bio: this.addDesignerForm.get('bio')?.value ?? '' }
      })
      .subscribe((value) => this.loadData(value));
  }

  updateNewDesigner() {
    this.apiClient
      .removeFontDesigner({ fontName: this.fontName, designerName: this.activeDesigner?.name ?? '' })
      .subscribe(() => {
        this.apiClient
          .addFontDesigner({
            fontName: this.fontName,
            body: { name: this.activeDesigner?.name ?? '', bio: this.editDesignerForm.get('bio')?.value ?? '' }
          })
          .subscribe((value) => {
            this.loadData(value);
            this.editDesignerOn = false;
          });
      });
  }

  removeDesigner() {
    this.apiClient
      .removeFontDesigner({ fontName: this.fontName, designerName: this.activeDesigner?.name ?? '' })
      .subscribe(() => this.loadData());
  }

  removeFile() {
    this.apiClient
      .removeFontFile({
        fontName: this.fontName,
        fontStyle: this.activeFile?.style ?? 'normal',
        fontWeight: this.activeFile?.weight ?? '400',
        fontType: this.activeFile?.type ?? 'woff2'
      })
      .subscribe(() => this.loadData());
  }

  openEditFile(file: Metadata) {
    this.activeFile = file;
    this.editFileOpen = true;
  }

  openDeleteFile(file: Metadata) {
    this.activeFile = file;
    this.removeFileOpen = true;
  }
}
