alter table `output_addresses_accumulate_in` add COLUMN `output_in_processed` smallint unsigned not null default 0;
alter table `output_addresses_accumulate_in` add COLUMN `output_out_processed` smallint unsigned not null default 0;
alter table `output_addresses_accumulate_out` add COLUMN `output_out_processed` smallint unsigned not null default 0;

create index `output_addresses_accumulate_in_output_processed` on `output_addresses_accumulate_in` (output_out_processed asc, output_in_processed asc, processed asc, created_at asc);
create index `output_addresses_accumulate_out_output_processed` on `output_addresses_accumulate_out` (output_out_processed asc, processed asc, created_at asc);

create index `output_addresses_accumulate_in_output_id` on `output_addresses_accumulate_in` (output_id asc);
create index `output_addresses_accumulate_out_output_id` on `output_addresses_accumulate_out` (output_id asc);
