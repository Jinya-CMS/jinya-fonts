/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { addFontDesigner } from '../fn/operations/add-font-designer';
import { AddFontDesigner$Params } from '../fn/operations/add-font-designer';
import { createFontFile } from '../fn/operations/create-font-file';
import { CreateFontFile$Params } from '../fn/operations/create-font-file';
import { createNewFont } from '../fn/operations/create-new-font';
import { CreateNewFont$Params } from '../fn/operations/create-new-font';
import { deleteFontByName } from '../fn/operations/delete-font-by-name';
import { DeleteFontByName$Params } from '../fn/operations/delete-font-by-name';
import { Designer } from '../models/designer';
import { getAllFonts } from '../fn/operations/get-all-fonts';
import { GetAllFonts$Params } from '../fn/operations/get-all-fonts';
import { getCustomFonts } from '../fn/operations/get-custom-fonts';
import { GetCustomFonts$Params } from '../fn/operations/get-custom-fonts';
import { getFontByName } from '../fn/operations/get-font-by-name';
import { GetFontByName$Params } from '../fn/operations/get-font-by-name';
import { getFontDesigners } from '../fn/operations/get-font-designers';
import { GetFontDesigners$Params } from '../fn/operations/get-font-designers';
import { getFontFiles } from '../fn/operations/get-font-files';
import { GetFontFiles$Params } from '../fn/operations/get-font-files';
import { getGoogleFonts } from '../fn/operations/get-google-fonts';
import { GetGoogleFonts$Params } from '../fn/operations/get-google-fonts';
import { getSettings } from '../fn/operations/get-settings';
import { GetSettings$Params } from '../fn/operations/get-settings';
import { Metadata } from '../models/metadata';
import { removeFontDesigner } from '../fn/operations/remove-font-designer';
import { RemoveFontDesigner$Params } from '../fn/operations/remove-font-designer';
import { removeFontFile } from '../fn/operations/remove-font-file';
import { RemoveFontFile$Params } from '../fn/operations/remove-font-file';
import { Settings } from '../models/settings';
import { updateFontByName } from '../fn/operations/update-font-by-name';
import { UpdateFontByName$Params } from '../fn/operations/update-font-by-name';
import { updateFontFile$Ttf } from '../fn/operations/update-font-file-ttf';
import { UpdateFontFile$Ttf$Params } from '../fn/operations/update-font-file-ttf';
import { updateFontFile$Woff2 } from '../fn/operations/update-font-file-woff-2';
import { UpdateFontFile$Woff2$Params } from '../fn/operations/update-font-file-woff-2';
import { updateSettings } from '../fn/operations/update-settings';
import { UpdateSettings$Params } from '../fn/operations/update-settings';
import { Webfont } from '../models/webfont';

@Injectable({ providedIn: 'root' })
export class ApiService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `getAllFonts()` */
  static readonly GetAllFontsPath = '/api/admin/font/all';

  /**
   * Gets all fonts.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getAllFonts()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllFonts$Response(params?: GetAllFonts$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Webfont>>> {
    return getAllFonts(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets all fonts.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getAllFonts$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllFonts(params?: GetAllFonts$Params, context?: HttpContext): Observable<Array<Webfont>> {
    return this.getAllFonts$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Webfont>>): Array<Webfont> => r.body)
    );
  }

  /** Path part for operation `createNewFont()` */
  static readonly CreateNewFontPath = '/api/admin/font/all';

  /**
   * Creates a new webfont.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createNewFont()` instead.
   *
   * This method doesn't expect any request body.
   */
  createNewFont$Response(params: CreateNewFont$Params, context?: HttpContext): Observable<StrictHttpResponse<Webfont>> {
    return createNewFont(this.http, this.rootUrl, params, context);
  }

  /**
   * Creates a new webfont.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createNewFont$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  createNewFont(params: CreateNewFont$Params, context?: HttpContext): Observable<Webfont> {
    return this.createNewFont$Response(params, context).pipe(
      map((r: StrictHttpResponse<Webfont>): Webfont => r.body)
    );
  }

  /** Path part for operation `getGoogleFonts()` */
  static readonly GetGoogleFontsPath = '/api/admin/font/google';

  /**
   * Gets all google fonts.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getGoogleFonts()` instead.
   *
   * This method doesn't expect any request body.
   */
  getGoogleFonts$Response(params?: GetGoogleFonts$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Webfont>>> {
    return getGoogleFonts(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets all google fonts.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getGoogleFonts$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getGoogleFonts(params?: GetGoogleFonts$Params, context?: HttpContext): Observable<Array<Webfont>> {
    return this.getGoogleFonts$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Webfont>>): Array<Webfont> => r.body)
    );
  }

  /** Path part for operation `getCustomFonts()` */
  static readonly GetCustomFontsPath = '/api/admin/font/custom';

  /**
   * Gets all custom fonts.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getCustomFonts()` instead.
   *
   * This method doesn't expect any request body.
   */
  getCustomFonts$Response(params?: GetCustomFonts$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Webfont>>> {
    return getCustomFonts(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets all custom fonts.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getCustomFonts$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getCustomFonts(params?: GetCustomFonts$Params, context?: HttpContext): Observable<Array<Webfont>> {
    return this.getCustomFonts$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Webfont>>): Array<Webfont> => r.body)
    );
  }

  /** Path part for operation `getFontByName()` */
  static readonly GetFontByNamePath = '/api/admin/font/{fontName}';

  /**
   * Gets the given font by name.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFontByName()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontByName$Response(params: GetFontByName$Params, context?: HttpContext): Observable<StrictHttpResponse<Webfont>> {
    return getFontByName(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the given font by name.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFontByName$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontByName(params: GetFontByName$Params, context?: HttpContext): Observable<Webfont> {
    return this.getFontByName$Response(params, context).pipe(
      map((r: StrictHttpResponse<Webfont>): Webfont => r.body)
    );
  }

  /** Path part for operation `updateFontByName()` */
  static readonly UpdateFontByNamePath = '/api/admin/font/{fontName}';

  /**
   * Updates the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updateFontByName()` instead.
   *
   * This method doesn't expect any request body.
   */
  updateFontByName$Response(params: UpdateFontByName$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updateFontByName(this.http, this.rootUrl, params, context);
  }

  /**
   * Updates the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updateFontByName$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  updateFontByName(params: UpdateFontByName$Params, context?: HttpContext): Observable<void> {
    return this.updateFontByName$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `deleteFontByName()` */
  static readonly DeleteFontByNamePath = '/api/admin/font/{fontName}';

  /**
   * Deletes the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `deleteFontByName()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteFontByName$Response(params: DeleteFontByName$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return deleteFontByName(this.http, this.rootUrl, params, context);
  }

  /**
   * Deletes the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `deleteFontByName$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteFontByName(params: DeleteFontByName$Params, context?: HttpContext): Observable<void> {
    return this.deleteFontByName$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `getFontFiles()` */
  static readonly GetFontFilesPath = '/api/admin/font/{fontName}/file';

  /**
   * Gets the files for the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFontFiles()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontFiles$Response(params: GetFontFiles$Params, context?: HttpContext): Observable<StrictHttpResponse<Metadata>> {
    return getFontFiles(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the files for the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFontFiles$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontFiles(params: GetFontFiles$Params, context?: HttpContext): Observable<Metadata> {
    return this.getFontFiles$Response(params, context).pipe(
      map((r: StrictHttpResponse<Metadata>): Metadata => r.body)
    );
  }

  /** Path part for operation `updateFontFile()` */
  static readonly UpdateFontFilePath = '/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}';

  /**
   * Replaces the font file of the given font and its parameters.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updateFontFile$Woff2()` instead.
   *
   * This method sends `font/woff2` and handles request body of type `font/woff2`.
   */
  updateFontFile$Woff2$Response(params: UpdateFontFile$Woff2$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updateFontFile$Woff2(this.http, this.rootUrl, params, context);
  }

  /**
   * Replaces the font file of the given font and its parameters.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updateFontFile$Woff2$Response()` instead.
   *
   * This method sends `font/woff2` and handles request body of type `font/woff2`.
   */
  updateFontFile$Woff2(params: UpdateFontFile$Woff2$Params, context?: HttpContext): Observable<void> {
    return this.updateFontFile$Woff2$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Replaces the font file of the given font and its parameters.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updateFontFile$Ttf()` instead.
   *
   * This method sends `font/ttf` and handles request body of type `font/ttf`.
   */
  updateFontFile$Ttf$Response(params: UpdateFontFile$Ttf$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updateFontFile$Ttf(this.http, this.rootUrl, params, context);
  }

  /**
   * Replaces the font file of the given font and its parameters.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updateFontFile$Ttf$Response()` instead.
   *
   * This method sends `font/ttf` and handles request body of type `font/ttf`.
   */
  updateFontFile$Ttf(params: UpdateFontFile$Ttf$Params, context?: HttpContext): Observable<void> {
    return this.updateFontFile$Ttf$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `createFontFile()` */
  static readonly CreateFontFilePath = '/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}';

  /**
   * Uploads a new font file with the given parameters.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createFontFile()` instead.
   *
   * This method sends `font/woff2` and handles request body of type `font/woff2`.
   */
  createFontFile$Response(params: CreateFontFile$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return createFontFile(this.http, this.rootUrl, params, context);
  }

  /**
   * Uploads a new font file with the given parameters.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createFontFile$Response()` instead.
   *
   * This method sends `font/woff2` and handles request body of type `font/woff2`.
   */
  createFontFile(params: CreateFontFile$Params, context?: HttpContext): Observable<void> {
    return this.createFontFile$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `removeFontFile()` */
  static readonly RemoveFontFilePath = '/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}';

  /**
   * Deletes the font file of the given font and its parameters.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `removeFontFile()` instead.
   *
   * This method doesn't expect any request body.
   */
  removeFontFile$Response(params: RemoveFontFile$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return removeFontFile(this.http, this.rootUrl, params, context);
  }

  /**
   * Deletes the font file of the given font and its parameters.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `removeFontFile$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  removeFontFile(params: RemoveFontFile$Params, context?: HttpContext): Observable<void> {
    return this.removeFontFile$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `getFontDesigners()` */
  static readonly GetFontDesignersPath = '/api/admin/font/{fontName}/designer';

  /**
   * Gets the designers for the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getFontDesigners()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontDesigners$Response(params: GetFontDesigners$Params, context?: HttpContext): Observable<StrictHttpResponse<Designer>> {
    return getFontDesigners(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the designers for the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getFontDesigners$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getFontDesigners(params: GetFontDesigners$Params, context?: HttpContext): Observable<Designer> {
    return this.getFontDesigners$Response(params, context).pipe(
      map((r: StrictHttpResponse<Designer>): Designer => r.body)
    );
  }

  /** Path part for operation `addFontDesigner()` */
  static readonly AddFontDesignerPath = '/api/admin/font/{fontName}/designer';

  /**
   * Adds a new designer to the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `addFontDesigner()` instead.
   *
   * This method doesn't expect any request body.
   */
  addFontDesigner$Response(params: AddFontDesigner$Params, context?: HttpContext): Observable<StrictHttpResponse<Designer>> {
    return addFontDesigner(this.http, this.rootUrl, params, context);
  }

  /**
   * Adds a new designer to the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `addFontDesigner$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  addFontDesigner(params: AddFontDesigner$Params, context?: HttpContext): Observable<Designer> {
    return this.addFontDesigner$Response(params, context).pipe(
      map((r: StrictHttpResponse<Designer>): Designer => r.body)
    );
  }

  /** Path part for operation `removeFontDesigner()` */
  static readonly RemoveFontDesignerPath = '/api/admin/font/{fontName}/designer/{designerName}';

  /**
   * Deletes given designer from the given font.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `removeFontDesigner()` instead.
   *
   * This method doesn't expect any request body.
   */
  removeFontDesigner$Response(params: RemoveFontDesigner$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return removeFontDesigner(this.http, this.rootUrl, params, context);
  }

  /**
   * Deletes given designer from the given font.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `removeFontDesigner$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  removeFontDesigner(params: RemoveFontDesigner$Params, context?: HttpContext): Observable<void> {
    return this.removeFontDesigner$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `getSettings()` */
  static readonly GetSettingsPath = '/api/admin/settings';

  /**
   * Gets the current settings.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getSettings()` instead.
   *
   * This method doesn't expect any request body.
   */
  getSettings$Response(params?: GetSettings$Params, context?: HttpContext): Observable<StrictHttpResponse<Settings>> {
    return getSettings(this.http, this.rootUrl, params, context);
  }

  /**
   * Gets the current settings.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getSettings$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getSettings(params?: GetSettings$Params, context?: HttpContext): Observable<Settings> {
    return this.getSettings$Response(params, context).pipe(
      map((r: StrictHttpResponse<Settings>): Settings => r.body)
    );
  }

  /** Path part for operation `updateSettings()` */
  static readonly UpdateSettingsPath = '/api/admin/settings';

  /**
   * Updates the current settings.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updateSettings()` instead.
   *
   * This method doesn't expect any request body.
   */
  updateSettings$Response(params: UpdateSettings$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updateSettings(this.http, this.rootUrl, params, context);
  }

  /**
   * Updates the current settings.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updateSettings$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  updateSettings(params: UpdateSettings$Params, context?: HttpContext): Observable<void> {
    return this.updateSettings$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

}
