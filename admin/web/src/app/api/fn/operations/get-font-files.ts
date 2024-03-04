/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Metadata } from '../../models/metadata';

export interface GetFontFiles$Params {
  fontName: string;
}

export function getFontFiles(
  http: HttpClient,
  rootUrl: string,
  params: GetFontFiles$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<Metadata>> {
  const rb = new RequestBuilder(rootUrl, getFontFiles.PATH, 'get');
  if (params) {
    rb.path('fontName', params.fontName, {});
  }

  return http.request(rb.build({ responseType: 'json', accept: 'application/json', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Metadata>;
    })
  );
}

getFontFiles.PATH = '/api/admin/font/{fontName}/file';
