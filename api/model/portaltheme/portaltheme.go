// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package portaltheme

// Background holds the background colors of a next-gen portal theme.
type Background struct {
	// Background color of the portal pages.
	// +kubebuilder:validation:Optional
	Page *string `json:"page,omitempty"`
	// Background color of the cards displayed on the portal pages.
	// +kubebuilder:validation:Optional
	Card *string `json:"card,omitempty"`
}

// Color is the color palette applied to the next-gen portal.
type Color struct {
	// +kubebuilder:validation:Optional
	Primary *string `json:"primary,omitempty"`
	// +kubebuilder:validation:Optional
	Secondary *string `json:"secondary,omitempty"`
	// +kubebuilder:validation:Optional
	Tertiary *string `json:"tertiary,omitempty"`
	// Color used to highlight errors on the portal.
	// +kubebuilder:validation:Optional
	Error *string `json:"error,omitempty"`
	// +kubebuilder:validation:Optional
	Background *Background `json:"background,omitempty"`
}

// Font is the typography applied to the next-gen portal.
type Font struct {
	// CSS font-family stack, e.g. `"Inter", sans-serif`. The console matches its own
	// list of stacks exactly, so a bare family name leaves its font selector empty.
	// +kubebuilder:validation:Optional
	FontFamily *string `json:"fontFamily,omitempty"`
}

// Definition is the next-gen portal theme definition. Leaving a section unset
// leaves the corresponding APIM defaults in place.
type Definition struct {
	// Custom CSS appended to the portal stylesheet.
	// +kubebuilder:validation:Optional
	CustomCss *string `json:"customCss,omitempty"`
	// +kubebuilder:validation:Optional
	Font *Font `json:"font,omitempty"`
	// +kubebuilder:validation:Optional
	Color *Color `json:"color,omitempty"`
}

// Type defines the specification of a PortalTheme resource.
// Images are inline data URIs, as the console produces on upload.
type Type struct {
	// Display name of the theme.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// The look and feel applied to the portal.
	// +kubebuilder:validation:Required
	Definition Definition `json:"definition"`
	// Logo displayed in the portal header.
	// +kubebuilder:validation:Optional
	Logo *string `json:"logo,omitempty"`
	// Alternate logo, used where the main logo does not fit.
	// +kubebuilder:validation:Optional
	OptionalLogo *string `json:"optionalLogo,omitempty"`
	// Icon displayed by the browser for the portal.
	// +kubebuilder:validation:Optional
	Favicon *string `json:"favicon,omitempty"`
	// Image displayed behind the portal header.
	// +kubebuilder:validation:Optional
	BackgroundImage *string `json:"backgroundImage,omitempty"`
}
