import { Component } from '@angular/core';
import { SettingsPanelComponent } from './settings-panel/settings-panel.component';
import { RouterLink, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [SettingsPanelComponent, RouterOutlet, RouterLink],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'settings-panel-fonts';
}
