/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { AddDesigner } from '../../models/add-designer';
import { Designer } from '../../models/designer';

export interface AddFontDesigner$Params {
  fontName: string;
  body?: AddDesigner;
}

export function addFontDesigner(
  http: HttpClient,
  rootUrl: string,
  params: AddFontDesigner$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<Designer>> {
  const rb = new RequestBuilder(rootUrl, addFontDesigner.PATH, 'post');
  if (params) {
    rb.path('fontName', params.fontName, {});
    rb.body(params.body, 'application/json');
  }

  return http.request(rb.build({ responseType: 'json', accept: 'application/json', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Designer>;
    })
  );
}

addFontDesigner.PATH = '/api/admin/font/{fontName}/designer';
