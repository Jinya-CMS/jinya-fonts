import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../api/services/api.service';
import { Status } from '../../api/models';
import { CheckCircle, XCircle } from 'lucide-angular';

@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrl: './status.component.scss',
  standalone: false
})
export class StatusComponent implements OnInit {
  loading = true;
  status!: Status;
  nextExecution: Date = new Date();

  constructor(private apiService: ApiService) {}

  ngOnInit() {
    this.apiService.getStatus().subscribe((status) => {
      this.status = status;
      if (status.jobNextExecution) {
        this.nextExecution = new Date(Date.parse(status.jobNextExecution));
      }

      this.loading = false;
    });
  }

  protected readonly CheckCircle = CheckCircle;
  protected readonly XCircle = XCircle;
}
