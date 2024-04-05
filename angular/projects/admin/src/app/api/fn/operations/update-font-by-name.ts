/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { UpdateWebfont } from '../../models/update-webfont';

export interface UpdateFontByName$Params {
  fontName: string;
  body?: UpdateWebfont;
}

export function updateFontByName(
  http: HttpClient,
  rootUrl: string,
  params: UpdateFontByName$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, updateFontByName.PATH, 'put');
  if (params) {
    rb.path('fontName', params.fontName, {});
    rb.body(params.body, 'application/json');
  }

  return http.request(rb.build({ responseType: 'text', accept: '*/*', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return (r as HttpResponse<any>).clone({ body: undefined }) as StrictHttpResponse<void>;
    })
  );
}

updateFontByName.PATH = '/api/admin/font/{fontName}';
