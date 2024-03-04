/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Webfont } from '../../models/webfont';

export interface GetAllFonts$Params {
}

export function getAllFonts(http: HttpClient, rootUrl: string, params?: GetAllFonts$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Webfont>>> {
  const rb = new RequestBuilder(rootUrl, getAllFonts.PATH, 'get');
  if (params) {
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Array<Webfont>>;
    })
  );
}

getAllFonts.PATH = '/api/admin/font/all';
