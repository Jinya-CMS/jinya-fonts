/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

export interface RemoveFontFile$Params {
  fontName: string;
  fontWeight: '100' | '200' | '300' | '400' | '500' | '600' | '700' | '800' | '900';
  fontStyle: 'normal' | 'italic';
  fontType: 'woff2' | 'ttf';
}

export function removeFontFile(
  http: HttpClient,
  rootUrl: string,
  params: RemoveFontFile$Params,
  context?: HttpContext
): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, removeFontFile.PATH, 'delete');
  if (params) {
    rb.path('fontName', params.fontName, {});
    rb.path('fontWeight', params.fontWeight, {});
    rb.path('fontStyle', params.fontStyle, {});
    rb.path('fontType', params.fontType, {});
  }

  return http.request(rb.build({ responseType: 'text', accept: '*/*', context })).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return (r as HttpResponse<any>).clone({ body: undefined }) as StrictHttpResponse<void>;
    })
  );
}

removeFontFile.PATH = '/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}';
