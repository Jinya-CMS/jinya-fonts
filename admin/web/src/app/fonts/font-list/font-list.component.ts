import { Component, OnInit } from '@angular/core';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { NgForOf, NgIf, NgSwitch, NgSwitchCase, NgSwitchDefault } from '@angular/common';
import { Check, LucideAngularModule, X } from 'lucide-angular';
import { RouterLink } from '@angular/router';
import { Designer } from '../../api/models/designer';
import { UiModule } from '../../ui/ui.module';

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
  standalone: true,
  imports: [NgForOf, LucideAngularModule, NgIf, RouterLink, NgSwitchCase, NgSwitchDefault, NgSwitch, UiModule],
  templateUrl: './font-list.component.html',
  styleUrl: './font-list.component.scss'
})
export class FontListComponent implements OnInit {
  fonts: Webfont[] = [];
  filteredFonts: Webfont[] = [];
  activeSideItem = ActiveSideItem.All;
  loading = true;

  activeFilter = {
    name: '',
    category: 'All'
  };

  protected readonly Check = Check;
  protected readonly X = X;
  protected readonly ActiveSideItem = ActiveSideItem;

  constructor(private apiClient: ApiService) {}

  ngOnInit(): void {
    this.getAllFonts();
  }

  getAllFonts(): void {
    this.loading = true;
    this.fonts = [];
    this.filteredFonts = [];
    this.apiClient.getAllFonts().subscribe((value) => {
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
    this.apiClient.getGoogleFonts().subscribe((value) => {
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
    this.apiClient.getCustomFonts().subscribe((value) => {
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
}
