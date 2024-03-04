/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Settings } from '../../models/settings';

export interface GetSettings$Params {
}

export function getSettings(http: HttpClient, rootUrl: string, params?: GetSettings$Params, context?: HttpContext): Observable<StrictHttpResponse<Settings>> {
  const rb = new RequestBuilder(rootUrl, getSettings.PATH, 'get');
  if (params) {
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Settings>;
    })
  );
}

getSettings.PATH = '/api/admin/settings';
