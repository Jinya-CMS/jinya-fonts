import { Alpine } from '../../../lib/alpine.js';
import { get } from '../../../lib/jinya-http.js';

import '../../lib/ui/online-indicator.js';

Alpine.data('statusData', () => ({
  jobIsScheduled: false,
  jobNextExecution: '',
  online: false,
  redisUp: false,
  mongoUp: false,
  servingWebsite: false,
  loading: true,
  async init() {
    const status = await get('/api/admin/status');
    this.jobIsScheduled = status.jobIsScheduled;
    this.jobNextExecution = new Date(Date.parse(status.jobNextExecution));
    this.online = status.online;
    this.redisUp = status.redisUp;
    this.servingWebsite = status.servingWebsite;
    this.loading = false;
  },
}));
