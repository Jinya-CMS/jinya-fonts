import { Component, computed, Input, OnInit } from '@angular/core';
import { Webfont } from '../api/models/webfont';
import { RouterLink } from '@angular/router';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';

@Component({
  selector: 'app-font-card',
  standalone: true,
  imports: [RouterLink],
  templateUrl: './font-card.component.html',
  styleUrl: './font-card.component.scss'
})
export class FontCardComponent implements OnInit {
  @Input() font!: Webfont;
  @Input() text!: string | null;

  designer = computed(() => this.font.designers?.map((designer) => designer.name).join(', '));

  styleUrl: SafeResourceUrl = '';

  constructor(private sanitizer: DomSanitizer) {}

  ngOnInit(): void {
    const url = `/css2?display=swap&family=${this.font.name}`;
    this.styleUrl = this.sanitizer.bypassSecurityTrustResourceUrl(url);
  }
}
