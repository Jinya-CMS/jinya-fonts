import { Component, ElementRef, HostListener, OnInit, ViewChild } from '@angular/core';
import { ApiService } from '../../api/services/api.service';
import { Webfont } from '../../api/models/webfont';
import { previewTexts } from '../../ui/preview-panel/preview-panel.component';
import { WebfontFilter } from '../../ui/filter-panel/filter-panel.component';

interface WebfontWithRow extends Webfont {
  rowStart: number;
  rowEnd: number;
  rowStartFull: number;
  rowEndFull: number;
}

@Component({
  selector: 'app-startpage',
  templateUrl: './start-page.component.html',
  styleUrl: './start-page.component.scss',
  standalone: false
})
export class StartPageComponent implements OnInit {
  constructor(private apiService: ApiService) {}

  webfonts: Webfont[] = [];
  filteredWebfonts: WebfontWithRow[] = [];
  fontsToShow: WebfontWithRow[] = [];
  loading = true;
  text = previewTexts.lorem;
  size = 24;

  @ViewChild('styleTester')
  styleTester!: ElementRef<HTMLDivElement>;

  protected readonly Math = Math;

  ngOnInit(): void {
    this.apiService.getFonts().subscribe((value) => {
      this.webfonts = value;
      this.loading = false;
      this.filter({
        display: true,
        handwriting: true,
        monospace: true,
        sansSerif: true,
        searchText: '',
        serif: true
      });
      this.setFontsToShow();
    });
  }

  filter(filters: WebfontFilter) {
    this.filteredWebfonts = this.webfonts
      .filter((font) => font.name.toLowerCase().includes(filters.searchText))
      .filter((item) => {
        let result = false;
        if (filters.sansSerif) {
          result = result || item.category.toLowerCase() === 'sans serif';
        }
        if (filters.serif) {
          result = result || item.category.toLowerCase() === 'serif';
        }
        if (filters.handwriting) {
          result = result || item.category.toLowerCase() === 'handwriting';
        }
        if (filters.display) {
          result = result || item.category.toLowerCase() === 'display';
        }
        if (filters.monospace) {
          result = result || item.category.toLowerCase() === 'monospace';
        }

        return result;
      })
      .map((font, idx) => ({
        ...font,
        rowStart: Math.floor((idx === 0 ? 0 : idx / 4) + 1),
        rowEnd: Math.floor((idx === 0 ? 0 : idx / 4) + 2),
        rowStartFull: idx + 1,
        rowEndFull: idx + 2
      }));
    this.setFontsToShow();
    window.scrollTo({
      top: 0
    });
  }

  setFontsToShow() {
    const isSmall = matchMedia('screen and (max-width: 768px)');

    let fontListHeight =
      window.innerHeight -
      (this.styleTester.nativeElement.clientHeight * 4.5 + this.styleTester.nativeElement.clientHeight * 6.5);
    if (isSmall) {
      fontListHeight = window.innerHeight;
    }

    const pageSize = Math.floor(fontListHeight / (this.styleTester.nativeElement.clientHeight * 13));

    let page = Math.floor((window.scrollY + window.innerHeight) / fontListHeight) - 1;
    if (page < 0) {
      page = 0;
    }

    let startIndex = (page - 2) * pageSize;
    if (startIndex < 0) {
      startIndex = 0;
    }

    const itemsPerLine = isSmall ? 1 : 4;
    const endIndex = page * pageSize * itemsPerLine + pageSize * 32;
    this.fontsToShow = this.filteredWebfonts.slice(startIndex * itemsPerLine, endIndex);
  }

  @HostListener('window:scroll')
  windowScroll() {
    this.setFontsToShow();
  }
}
