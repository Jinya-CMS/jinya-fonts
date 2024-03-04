import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { Observable } from 'rxjs';
import { Metadata } from '../../api/models/metadata';

@Component({
  selector: 'app-update-file-dialog',
  templateUrl: './update-file-dialog.component.html'
})
export class UpdateFileDialogComponent {
  editFileForm = new FormGroup({
    file: new FormControl<File | null>(null, { nonNullable: true })
  });

  @Input() open = false;
  @Input() font!: Webfont;
  @Input() file!: Metadata;

  @Output() saved = new EventEmitter();

  constructor(private apiClient: ApiService) {}

  updateFile() {
    if (this.editFileForm.invalid) {
      return;
    }

    const file = this.editFileForm.get('file')?.value;
    let saveObservable = new Observable<void>();
    if (file && this.file.type === 'woff2') {
      saveObservable = this.apiClient.updateFontFile$Woff2({
        body: file,
        fontName: this.font.name,
        fontStyle: this.file.style,
        fontWeight: this.file.weight,
        fontType: 'woff2'
      });
    } else if (file && this.file.type === 'ttf') {
      saveObservable = this.apiClient.updateFontFile$Ttf({
        body: file,
        fontName: this.font.name,
        fontStyle: this.file.style,
        fontWeight: this.file.weight,
        fontType: 'ttf'
      });
    }

    saveObservable.subscribe(() => {
      this.open = false;
      this.saved.emit();
    });
  }

  updateFilePicker($event: Event) {
    // @ts-expect-error The event cannot be null
    const file = ($event.target as HTMLInputElement).files[0];
    this.editFileForm.get('file')?.patchValue(file);
  }
}
