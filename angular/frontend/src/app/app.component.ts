import { Component } from '@angular/core';
import { SettingsPanelComponent } from './settingsPanel/settingsPanel.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [SettingsPanelComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'settingsPanel-fonts';
}
