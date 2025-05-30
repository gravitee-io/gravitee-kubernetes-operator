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

package logback

var Config = `
<configuration>

    <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
        <encoder class="ch.qos.logback.core.encoder.LayoutWrappingEncoder">
            <layout class="ch.qos.logback.contrib.json.classic.JsonLayout">
                <jsonFormatter
                        class="ch.qos.logback.contrib.jackson.JacksonJsonFormatter">
                </jsonFormatter>
                <appendLineSeparator>true</appendLineSeparator>
                <timestampFormat>yyyy-MM-dd'T'HH:mm:ss.SSSXX</timestampFormat>
            </layout>
        </encoder>
    </appender>

    <appender name="async-console" class="ch.qos.logback.classic.AsyncAppender">
        <appender-ref ref="STDOUT" />
    </appender>

    <logger name="io.gravitee" level="INFO" />
    <logger name="org.reflections" level="WARN" />
    <logger name="org.springframework" level="WARN" />
    <logger name="org.eclipse.jetty" level="WARN" />
    <logger name="com.graviteesource.reactor.nativekafka" level="DEBUG" />

    <root level="ERROR">
        <appender-ref ref="async-console" />
    </root>

</configuration>
`
