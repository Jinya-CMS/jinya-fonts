import { Component, Input, OnInit } from '@angular/core';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { UiModule } from '../../ui/ui.module';
import { Check, LucideAngularModule, X } from 'lucide-angular';
import { Designer } from '../../api/models/designer';
import { Metadata } from '../../api/models/metadata';

enum ActiveSideListItem {
  Details,
  Designers,
  Files
}

@Component({
  selector: 'app-font-details',
  standalone: true,
  imports: [UiModule, LucideAngularModule],
  templateUrl: './font-details.component.html',
  styleUrl: './font-details.component.scss'
})
export class FontDetailsComponent implements OnInit {
  @Input() fontName!: string;

  font!: Webfont;
  activeDesigner: Designer | null = null;
  loading = true;
  activeSideListItem = ActiveSideListItem.Details;

  protected readonly ActiveSideListItem = ActiveSideListItem;
  protected readonly Check = Check;
  protected readonly X = X;

  constructor(private apiClient: ApiService) {}

  ngOnInit(): void {
    this.apiClient.getFontByName({ fontName: this.fontName }).subscribe((value) => {
      this.font = value;
      this.loading = false;
      this.activeDesigner = this.font.designers?.at(0) ?? null;
    });
  }

  changeSideItem($event: MouseEvent, tab: ActiveSideListItem) {
    $event.preventDefault();
    this.activeSideListItem = tab;
  }

  selectDesigner($event: MouseEvent, designer: Designer) {
    $event.preventDefault();
    this.activeDesigner = designer;
  }
}
