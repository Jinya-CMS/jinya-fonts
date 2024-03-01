import { Component } from '@angular/core';
import { FontCard } from '../fontCard';
import { SearchbarComponent } from '../searchbar/searchbar.component';

@Component({
  selector: 'app-settingsPanel',
  standalone: true,
  imports: [SearchbarComponent],
  templateUrl: './settingsPanel.component.html',
  styleUrl: './settingsPanel.component.scss'
})
export class SettingsPanelComponent {
  fontList: FontCard[] = [];
  filteredFontList: FontCard[] = [];
}
