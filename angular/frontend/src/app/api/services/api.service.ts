/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { apiFontFontNameDesignerGet } from '../fn/operations/api-font-font-name-designer-get';
import { ApiFontFontNameDesignerGet$Params } from '../fn/operations/api-font-font-name-designer-get';
import { apiFontFontNameFileGet } from '../fn/operations/api-font-font-name-file-get';
import { ApiFontFontNameFileGet$Params } from '../fn/operations/api-font-font-name-file-get';
import { apiFontFontNameGet } from '../fn/operations/api-font-font-name-get';
import { ApiFontFontNameGet$Params } from '../fn/operations/api-font-font-name-get';
import { apiFontGet } from '../fn/operations/api-font-get';
import { ApiFontGet$Params } from '../fn/operations/api-font-get';
import { Designer } from '../models/designer';
import { Metadata } from '../models/metadata';
import { Webfont } from '../models/webfont';

@Injectable({ providedIn: 'root' })
export class ApiService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `apiFontGet()` */
  static readonly ApiFontGetPath = '/api/font';

  /**
   * Gets all fonts.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `apiFontGet()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontGet$Response(
    params?: ApiFontGet$Params,
    context?: HttpContext
  ): Observable<StrictHttpResponse<Array<Webfont>>> {
    return apiFontGet(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets all fonts.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `apiFontGet$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontGet(params?: ApiFontGet$Params, context?: HttpContext): Observable<Array<Webfont>> {
    return this.apiFontGet$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Webfont>>): Array<Webfont> => r.body)
    );
  }

  /** Path part for operation `apiFontFontNameGet()` */
  static readonly ApiFontFontNameGetPath = '/api/font/{fontName}';

  /**
   * Gets the given font by name.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `apiFontFontNameGet()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontFontNameGet$Response(
    params: ApiFontFontNameGet$Params,
    context?: HttpContext
  ): Observable<StrictHttpResponse<Webfont>> {
    return apiFontFontNameGet(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the given font by name.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `apiFontFontNameGet$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontFontNameGet(params: ApiFontFontNameGet$Params, context?: HttpContext): Observable<Webfont> {
    return this.apiFontFontNameGet$Response(params, context).pipe(
      map((r: StrictHttpResponse<Webfont>): Webfont => r.body)
    );
  }

  /** Path part for operation `apiFontFontNameFileGet()` */
  static readonly ApiFontFontNameFileGetPath = '/api/font/{fontName}/file';

  /**
   * Gets the files for the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `apiFontFontNameFileGet()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontFontNameFileGet$Response(
    params: ApiFontFontNameFileGet$Params,
    context?: HttpContext
  ): Observable<StrictHttpResponse<Metadata>> {
    return apiFontFontNameFileGet(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the files for the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `apiFontFontNameFileGet$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontFontNameFileGet(params: ApiFontFontNameFileGet$Params, context?: HttpContext): Observable<Metadata> {
    return this.apiFontFontNameFileGet$Response(params, context).pipe(
      map((r: StrictHttpResponse<Metadata>): Metadata => r.body)
    );
  }

  /** Path part for operation `apiFontFontNameDesignerGet()` */
  static readonly ApiFontFontNameDesignerGetPath = '/api/font/{fontName}/designer';

  /**
   * Gets the designers for the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `apiFontFontNameDesignerGet()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontFontNameDesignerGet$Response(
    params: ApiFontFontNameDesignerGet$Params,
    context?: HttpContext
  ): Observable<StrictHttpResponse<Designer>> {
    return apiFontFontNameDesignerGet(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the designers for the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `apiFontFontNameDesignerGet$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  apiFontFontNameDesignerGet(params: ApiFontFontNameDesignerGet$Params, context?: HttpContext): Observable<Designer> {
    return this.apiFontFontNameDesignerGet$Response(params, context).pipe(
      map((r: StrictHttpResponse<Designer>): Designer => r.body)
    );
  }
}
