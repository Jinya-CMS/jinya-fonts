/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Webfont } from '../../models/webfont';

export interface GetFonts$Params {}

export function getFonts(
  http: HttpClient,
  rootUrl: string,
  params?: GetFonts$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<Array<Webfont>>> {
  const rb = new RequestBuilder(rootUrl, getFonts.PATH, 'get');
  if (params) {
  }

  return http.request(rb.build({ responseType: 'json', accept: 'application/json', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Array<Webfont>>;
    })
  );
}

getFonts.PATH = '/api/font';
