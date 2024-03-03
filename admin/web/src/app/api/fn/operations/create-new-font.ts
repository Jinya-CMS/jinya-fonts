/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Webfont } from '../../models/webfont';

export interface CreateNewFont$Params {
}

export function createNewFont(http: HttpClient, rootUrl: string, params: CreateNewFont$Params, context?: HttpContext): Observable<StrictHttpResponse<Webfont>> {
  const rb = new RequestBuilder(rootUrl, createNewFont.PATH, 'post');
  if (params) {
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Webfont>;
    })
  );
}

createNewFont.PATH = '/api/admin/font/all';
