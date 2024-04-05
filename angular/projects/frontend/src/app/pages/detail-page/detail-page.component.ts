import { Component, OnInit } from '@angular/core';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { ActivatedRoute } from '@angular/router';
import { Metadata } from '../../api/models/metadata';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { previewTexts } from '../../ui/preview-panel/preview-panel.component';

@Component({
  selector: 'app-detailpage',
  templateUrl: './detail-page.component.html',
  styleUrl: './detail-page.component.scss'
})
export class DetailPageComponent implements OnInit {
  loading = true;
  font!: Webfont;

  designers = '';
  fontHtml = '';
  fontCss = '';
  previewText = previewTexts.lorem;
  styleUrl: SafeResourceUrl = '';

  selectedFiles: Metadata[] = [];
  previewSize = 24;

  constructor(
    private apiService: ApiService,
    private route: ActivatedRoute,
    private sanitizer: DomSanitizer
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe(({ font }) => {
      this.loading = true;
      this.apiService.getFontByName({ fontName: font }).subscribe((value) => {
        this.font = value;
        this.font.fonts?.sort((a, b) => {
          const w = a.weight.localeCompare(b.weight);
          if (w === 0) {
            return a.style.localeCompare(b.style) * -1;
          }

          return w;
        });
        this.designers = value.designers?.map((designer) => designer.name).join(', ') ?? '';
        this.loading = false;
        const url = `/css2?display=swap&family=${this.font.name}:all`;
        this.styleUrl = this.sanitizer.bypassSecurityTrustResourceUrl(url);
        this.fontCss = `body {
    font-family: "${this.font.name}", ${this.font.category.toLowerCase().replace(' ', '-')};
}`;
        this.updateHtml();
      });
    });
  }

  toggleFont(file: Metadata) {
    if (this.selectedFiles.includes(file)) {
      this.selectedFiles = this.selectedFiles.filter((f) => f !== file);
    } else {
      this.selectedFiles.push(file);
    }

    this.selectedFiles = this.selectedFiles.sort((a, b) => {
      const w = a.weight.localeCompare(b.weight);
      if (w === 0) {
        return a.style.localeCompare(b.style) * -1;
      }

      return w;
    });

    this.updateHtml();
  }

  updateHtml() {
    let params = '';
    if (this.selectedFiles.length > 0) {
      params = `&ital,wght@${this.selectedFiles.map((item) => `${item.style === 'italic' ? 1 : 0},${item.weight}`).join('%3B')}`;
    }

    this.fontHtml = `<link rel="stylesheet" type="text/css" href="${location.origin}/css2?family=${this.font.name}${params}">`;
  }
}
