import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../api/services/api.service';
import { FormControl, FormGroup } from '@angular/forms';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss',
  standalone: false
})
export class SettingsComponent implements OnInit {
  onlyFetchFonts: string[] = [];

  loading = true;
  loaded = false;

  settingsForm = new FormGroup({
    syncInterval: new FormControl(''),
    syncEnabled: new FormControl(false)
  });
  addFilterForm = new FormGroup({
    filterName: new FormControl('')
  });

  constructor(private apiService: ApiService) {}

  ngOnInit() {
    this.loadData();
  }

  loadData() {
    if (!this.loaded) {
      this.loading = true;
    }

    this.apiService.getSettings().subscribe((value) => {
      this.onlyFetchFonts = value.filterByName;
      this.settingsForm.get('syncInterval')?.setValue(value.syncInterval ?? '0 0 1 * *');
      this.settingsForm.get('syncEnabled')?.setValue(value.syncEnabled ?? true);
      this.loading = false;
      this.loaded = true;
    });
  }

  addFilter() {
    if (this.addFilterForm.invalid) {
      return;
    }

    const filterName = this.addFilterForm.get('filterName')?.value ?? '';
    if (filterName && !this.onlyFetchFonts.includes(filterName)) {
      this.onlyFetchFonts.push(this.addFilterForm.get('filterName')?.value ?? '');
      this.addFilterForm.reset();
    }
  }

  saveSettings() {
    this.apiService
      .updateSettings({
        body: {
          syncInterval: this.settingsForm.get('syncInterval')?.value ?? '0 0 1 * *',
          syncEnabled: this.settingsForm.get('syncEnabled')?.value ?? true,
          filterByName: this.onlyFetchFonts
        }
      })
      .subscribe(() => this.loadData());
  }

  removeFilteredFont(font: string) {
    this.onlyFetchFonts = this.onlyFetchFonts.filter((value) => value !== font);
  }
}
