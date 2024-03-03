/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';


export interface RemoveFontDesigner$Params {
  fontName: string;
  designerName: string;
}

export function removeFontDesigner(http: HttpClient, rootUrl: string, params: RemoveFontDesigner$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, removeFontDesigner.PATH, 'delete');
  if (params) {
    rb.path('fontName', params.fontName, {});
    rb.path('designerName', params.designerName, {});
  }

  return http.request(
    rb.build({ responseType: 'text', accept: '*/*', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return (r as HttpResponse<any>).clone({ body: undefined }) as StrictHttpResponse<void>;
    })
  );
}

removeFontDesigner.PATH = '/api/admin/font/{fontName}/designer/{designerName}';
