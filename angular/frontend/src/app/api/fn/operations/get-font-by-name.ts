/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Webfont } from '../../models/webfont';

export interface GetFontByName$Params {
  fontName: string;
}

export function getFontByName(
  http: HttpClient,
  rootUrl: string,
  params: GetFontByName$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<Webfont>> {
  const rb = new RequestBuilder(rootUrl, getFontByName.PATH, 'get');
  if (params) {
    rb.path('fontName', params.fontName, {});
  }

  return http.request(rb.build({ responseType: 'json', accept: 'application/json', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Webfont>;
    })
  );
}

getFontByName.PATH = '/api/font/{fontName}';
