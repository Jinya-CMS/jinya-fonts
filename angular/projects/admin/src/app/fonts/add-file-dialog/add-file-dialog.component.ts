import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { Webfont } from '../../api/models/webfont';
import { ApiService } from '../../api/services/api.service';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-add-file-dialog',
  templateUrl: './add-file-dialog.component.html'
})
export class AddFileDialogComponent {
  createFileForm = new FormGroup({
    weight: new FormControl<'100' | '200' | '300' | '400' | '500' | '600' | '700' | '800' | '900'>('400', {
      nonNullable: true
    }),
    style: new FormControl<'normal' | 'italic'>('normal', { nonNullable: true }),
    file: new FormControl<File | null>(null, { nonNullable: true })
  });

  @Input() open = false;
  @Input() font!: Webfont;

  @Output() saved = new EventEmitter();

  constructor(private apiService: ApiService) {}

  createFile() {
    if (this.createFileForm.invalid) {
      return;
    }

    const file = this.createFileForm.get('file')?.value;
    let saveObservable = new Observable<void>();
    if (file && file.type === 'font/woff2') {
      saveObservable = this.apiService.createFontFile$Woff2({
        body: file,
        fontName: this.font.name,
        fontStyle: this.createFileForm.get('style')?.value ?? 'normal',
        fontWeight: this.createFileForm.get('weight')?.value ?? '400',
        fontType: 'woff2'
      });
    } else if (file && file.type === 'font/ttf') {
      saveObservable = this.apiService.createFontFile$Ttf({
        body: file,
        fontName: this.font.name,
        fontStyle: this.createFileForm.get('style')?.value ?? 'normal',
        fontWeight: this.createFileForm.get('weight')?.value ?? '400',
        fontType: 'ttf'
      });
    }

    saveObservable.subscribe(() => {
      this.open = false;
      this.saved.emit();
    });
  }

  updateFile($event: Event) {
    // @ts-expect-error The event cannot be null
    const file = ($event.target as HTMLInputElement).files[0];
    this.createFileForm.get('file')?.patchValue(file);
  }
}
