/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Designer } from '../../models/designer';

export interface GetFontDesigners$Params {
  fontName: string;
}

export function getFontDesigners(
  http: HttpClient,
  rootUrl: string,
  params: GetFontDesigners$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<Designer>> {
  const rb = new RequestBuilder(rootUrl, getFontDesigners.PATH, 'get');
  if (params) {
    rb.path('fontName', params.fontName, {});
  }

  return http.request(rb.build({ responseType: 'json', accept: 'application/json', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Designer>;
    })
  );
}

getFontDesigners.PATH = '/api/admin/font/{fontName}/designer';
