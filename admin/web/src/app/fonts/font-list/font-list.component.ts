import { Component, Input, OnInit } from '@angular/core';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { Check, X } from 'lucide-angular';
import { Designer } from '../../api/models/designer';
import { FontSyncEvents, FontSyncService } from '../font-sync.service';

interface FontFilter {
  category: string;
  name: string;
}

enum ActiveSideItem {
  All,
  Google,
  Custom
}

@Component({
  selector: 'app-font-list',
  templateUrl: './font-list.component.html',
  styleUrl: './font-list.component.scss'
})
export class FontListComponent implements OnInit {
  @Input() type!: string;

  fonts: Webfont[] = [];
  filteredFonts: Webfont[] = [];
  activeSideItem = ActiveSideItem.All;
  loading = true;
  deleteFontOpen = false;
  selectedFont: Webfont | null = null;

  activeFilter = {
    name: '',
    category: 'All'
  };

  protected readonly Check = Check;
  protected readonly X = X;
  protected readonly ActiveSideItem = ActiveSideItem;

  constructor(
    private apiService: ApiService,
    private fontSyncService: FontSyncService
  ) {
    fontSyncService.subscribe(FontSyncEvents.Success).subscribe(() => {
      switch (this.activeSideItem) {
        case ActiveSideItem.All:
          this.getAllFonts();
          break;
        case ActiveSideItem.Google:
          this.getGoogleFonts();
          break;
      }
    });
  }

  ngOnInit(): void {
    switch (this.type.toLowerCase()) {
      case 'google':
        this.activeSideItem = ActiveSideItem.Google;
        this.getGoogleFonts();
        break;
      case 'custom':
        this.activeSideItem = ActiveSideItem.Custom;
        this.getCustomFonts();
        break;
      default:
        this.activeSideItem = ActiveSideItem.All;
        this.getAllFonts();
        break;
    }
  }

  getAllFonts(): void {
    this.loading = true;
    this.fonts = [];
    this.filteredFonts = [];
    this.apiService.getAllFonts().subscribe((value) => {
      this.filteredFonts = value;
      this.fonts = value;
      this.loading = false;
    });
    this.activeFilter = {
      name: '',
      category: 'All'
    };
  }

  getGoogleFonts(): void {
    this.loading = true;
    this.fonts = [];
    this.filteredFonts = [];
    this.apiService.getGoogleFonts().subscribe((value) => {
      this.filteredFonts = value;
      this.fonts = value;
      this.loading = false;
    });
    this.activeFilter = {
      name: '',
      category: 'All'
    };
  }

  getCustomFonts(): void {
    this.loading = true;
    this.fonts = [];
    this.filteredFonts = [];
    this.apiService.getCustomFonts().subscribe((value) => {
      this.filteredFonts = value;
      this.fonts = value;
      this.loading = false;
    });
    this.activeFilter = {
      name: '',
      category: 'All'
    };
  }

  designersToString(designers: Array<Designer> | undefined) {
    return designers?.map((d) => d.name).join(', ') ?? '';
  }

  filterFonts(filter: FontFilter) {
    this.filteredFonts = this.fonts.filter((f) => f.name.toLowerCase().includes(filter.name.toLowerCase()));
    if (filter.category !== 'All') {
      this.filteredFonts = this.filteredFonts.filter((f) => f.category === filter.category);
    }

    this.activeFilter = filter;
  }

  changeSideItem($event: MouseEvent, item: ActiveSideItem) {
    $event.preventDefault();
    this.activeSideItem = item;
    switch (item) {
      case ActiveSideItem.All:
        this.getAllFonts();
        break;
      case ActiveSideItem.Google:
        this.getGoogleFonts();
        break;
      case ActiveSideItem.Custom:
        this.getCustomFonts();
        break;
    }
  }

  openDelete(font: Webfont) {
    this.deleteFontOpen = true;
    this.selectedFont = font;
  }

  deleteFont() {
    this.apiService.deleteFontByName({ fontName: this.selectedFont?.name ?? '' }).subscribe(() => {
      this.deleteFontOpen = false;
      this.fonts = this.fonts.filter((font) => font.name !== this.selectedFont?.name);
      this.filterFonts(this.activeFilter);
    });
  }

  syncFonts() {
    this.fontSyncService.syncFonts();
  }
}
