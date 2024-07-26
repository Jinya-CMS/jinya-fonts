/* tslint:disable */
/* eslint-disable */
export interface Status {
  jobIsScheduled: boolean;
  jobNextExecution?: string;
  mongoUp: boolean;
  online: boolean;
  redisUp: boolean;
  servingWebsite: boolean;
}
